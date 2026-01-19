package clean

import (
	"context"
	"errors"
	"testing"

	"github.com/pahluwalia-tcloud/together-kubelogin/mocks/github.com/pahluwalia-tcloud/together-kubelogin/pkg/tokencache/repository_mock"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/testing/logger"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/tokencache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClean_Do(t *testing.T) {
	t.Run("Success_BothStorages", func(t *testing.T) {
		mockRepo := repository_mock.NewMockInterface(t)
		testLogger := logger.New(t)

		// Expect disk deletion
		mockRepo.EXPECT().DeleteAll(tokencache.Config{
			Directory: "/test/cache/dir",
			Storage:   tokencache.StorageDisk,
		}).Return(nil).Once()

		// Expect keyring deletion
		mockRepo.EXPECT().DeleteAll(tokencache.Config{
			Directory: "/test/cache/dir",
			Storage:   tokencache.StorageKeyring,
		}).Return(nil).Once()

		clean := Clean{
			TokenCacheRepository: mockRepo,
			Logger:               testLogger,
		}

		err := clean.Do(context.Background(), Input{
			TokenCacheDir: "/test/cache/dir",
		})
		require.NoError(t, err)
	})

	t.Run("DiskDeletionError", func(t *testing.T) {
		mockRepo := repository_mock.NewMockInterface(t)
		testLogger := logger.New(t)

		diskError := errors.New("disk deletion failed")
		mockRepo.EXPECT().DeleteAll(tokencache.Config{
			Directory: "/test/cache/dir",
			Storage:   tokencache.StorageDisk,
		}).Return(diskError).Once()

		clean := Clean{
			TokenCacheRepository: mockRepo,
			Logger:               testLogger,
		}

		err := clean.Do(context.Background(), Input{
			TokenCacheDir: "/test/cache/dir",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete the token cache from /test/cache/dir")
		assert.Contains(t, err.Error(), "disk deletion failed")
	})

	t.Run("KeyringDeletionError_NotFatal", func(t *testing.T) {
		mockRepo := repository_mock.NewMockInterface(t)
		testLogger := logger.New(t)

		// Disk deletion succeeds
		mockRepo.EXPECT().DeleteAll(tokencache.Config{
			Directory: "/test/cache/dir",
			Storage:   tokencache.StorageDisk,
		}).Return(nil).Once()

		// Keyring deletion fails but should not stop execution
		keyringError := errors.New("keyring not available")
		mockRepo.EXPECT().DeleteAll(tokencache.Config{
			Directory: "/test/cache/dir",
			Storage:   tokencache.StorageKeyring,
		}).Return(keyringError).Once()

		clean := Clean{
			TokenCacheRepository: mockRepo,
			Logger:               testLogger,
		}

		err := clean.Do(context.Background(), Input{
			TokenCacheDir: "/test/cache/dir",
		})
		require.NoError(t, err, "keyring error should not be fatal")
	})

	t.Run("EmptyTokenCacheDir", func(t *testing.T) {
		mockRepo := repository_mock.NewMockInterface(t)
		testLogger := logger.New(t)

		// Should still attempt to delete with empty directory
		mockRepo.EXPECT().DeleteAll(tokencache.Config{
			Directory: "",
			Storage:   tokencache.StorageDisk,
		}).Return(nil).Once()

		mockRepo.EXPECT().DeleteAll(tokencache.Config{
			Directory: "",
			Storage:   tokencache.StorageKeyring,
		}).Return(nil).Once()

		clean := Clean{
			TokenCacheRepository: mockRepo,
			Logger:               testLogger,
		}

		err := clean.Do(context.Background(), Input{
			TokenCacheDir: "",
		})
		require.NoError(t, err)
	})
}
