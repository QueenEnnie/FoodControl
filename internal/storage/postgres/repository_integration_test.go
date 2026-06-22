//go:build integration

package postgres

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"food-control/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const integrationTestTimeout = 2 * time.Minute

func TestRepositoryIntegration(t *testing.T) {
	repo, pool := startPostgres(t)

	t.Run("create and get product", func(t *testing.T) {
		resetDatabase(t, pool)
		expiresAt := time.Now().UTC().AddDate(0, 0, 7).Truncate(24 * time.Hour)
		created, err := repo.CreateProduct(t.Context(), domain.CreateProductRequest{
			Name: "milk", Description: "integration test product", Quantity: 1.5,
			Unit: "l", ExpirationDate: &expiresAt,
		})
		if err != nil {
			t.Fatalf("CreateProduct() error = %v", err)
		}

		got, err := repo.GetProduct(t.Context(), created.ID)
		if err != nil {
			t.Fatalf("GetProduct() error = %v", err)
		}
		if got.Name != "milk" || got.Description != "integration test product" || got.Unit != "l" {
			t.Fatalf("GetProduct() = %+v, want created product", got)
		}
		if math.Abs(got.Quantity-1.5) > 0.000001 {
			t.Errorf("GetProduct().Quantity = %v, want 1.5", got.Quantity)
		}
	})

	t.Run("recipe availability excludes expired stock", func(t *testing.T) {
		resetDatabase(t, pool)
		today := time.Now().UTC().Truncate(24 * time.Hour)
		freshUntil := today.AddDate(0, 0, 7)
		expiredAt := today.AddDate(0, 0, -1)

		for _, input := range []domain.CreateProductRequest{
			{Name: "milk", Quantity: 0.3, Unit: "l", ExpirationDate: &freshUntil},
			{Name: "MILK", Quantity: 1, Unit: "L", ExpirationDate: &expiredAt},
		} {
			if _, err := repo.CreateProduct(t.Context(), input); err != nil {
				t.Fatalf("CreateProduct(%q) error = %v", input.Name, err)
			}
		}

		recipe, err := repo.CreateRecipe(t.Context(), domain.CreateRecipeRequest{
			Name: "Pancakes",
			Ingredients: []domain.RecipeIngredient{
				{ProductName: "Milk", RequiredQuantity: 0.5, Unit: "l"},
			},
		})
		if err != nil {
			t.Fatalf("CreateRecipe() error = %v", err)
		}

		availability, err := repo.GetRecipeAvailability(t.Context(), recipe.ID)
		if err != nil {
			t.Fatalf("GetRecipeAvailability() error = %v", err)
		}
		if availability.CanCook {
			t.Error("GetRecipeAvailability().CanCook = true, want false")
		}
		if availability.MissingCount != 1 {
			t.Errorf("MissingCount = %d, want 1", availability.MissingCount)
		}
		if len(availability.Ingredients) != 1 {
			t.Fatalf("len(Ingredients) = %d, want 1", len(availability.Ingredients))
		}
		if got := availability.Ingredients[0].AvailableQuantity; math.Abs(got-0.3) > 0.000001 {
			t.Errorf("available quantity = %v, want 0.3; expired stock must be ignored", got)
		}
	})
}

func startPostgres(t *testing.T) (*Repository, *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), integrationTestTimeout)
	t.Cleanup(cancel)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:16-alpine", ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB": "food_control_test", "POSTGRES_USER": "test", "POSTGRES_PASSWORD": "test",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(integrationTestTimeout),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start PostgreSQL container: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		if err := container.Terminate(cleanupCtx); err != nil {
			t.Errorf("terminate PostgreSQL container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("get container port: %v", err)
	}
	connectionString := fmt.Sprintf(
		"postgres://test:test@%s:%s/food_control_test?sslmode=disable", host, port.Port(),
	)
	pool, err := NewPool(ctx, connectionString)
	if err != nil {
		t.Fatalf("connect to PostgreSQL container: %v", err)
	}
	t.Cleanup(pool.Close)

	applyMigrations(t, pool)
	return New(pool), pool
}

func applyMigrations(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve integration test path")
	}
	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
	for _, name := range []string{"001_init.sql", "002_products_soft_delete.sql"} {
		migration, err := os.ReadFile(filepath.Join(projectRoot, "migrations", name))
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if _, err := pool.Exec(t.Context(), string(migration)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
}

func resetDatabase(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(t.Context(), `
		TRUNCATE TABLE recipe_ingredients, recipes, products
		RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("reset integration database: %v", err)
	}
}
