package notifiers

import (
	"github.com/statping/statping/database"
	"github.com/statping/statping/types/notifications"
	"github.com/statping/statping/types/null"
	"github.com/statping/statping/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	DISCORD_URL string
)

func TestDiscordNotifier(t *testing.T) {
	err := utils.InitLogs()
	require.Nil(t, err)
	DISCORD_URL = utils.Params.GetString("DISCORD_URL")

	db, err := database.OpenTester()
	require.Nil(t, err)
	db.AutoMigrate(&notifications.Notification{})
	notifications.SetDB(db)

	if DISCORD_URL == "" {
		t.Log("discord notifier testing skipped, missing DISCORD_URL environment variable")
		t.SkipNow()
	}

	t.Run("Load discord", func(t *testing.T) {
		Discorder.Host = DISCORD_URL
		Discorder.Delay = time.Duration(100 * time.Millisecond)
		Discorder.Enabled = null.NewNullBool(true)

		Add(Discorder)

		assert.Equal(t, "Hunter Long", Discorder.Author)
		assert.Equal(t, DISCORD_URL, Discorder.Host)
	})

	t.Run("discord Notifier Tester", func(t *testing.T) {
		assert.True(t, Discorder.CanSend())
	})

	t.Run("discord OnFailure", func(t *testing.T) {
		_, err := Discorder.OnFailure(exampleService, exampleFailure)
		assert.Nil(t, err)
	})

	t.Run("discord OnSuccess", func(t *testing.T) {
		_, err := Discorder.OnSuccess(exampleService)
		assert.Nil(t, err)
	})

	t.Run("discord Test", func(t *testing.T) {
		_, err := Discorder.OnTest()
		assert.Nil(t, err)
	})

}
