package goblar

// Run is the zero-brain entrypoint for go-blar.
// It creates an App, registers models, and starts the server.
func Run(models ...any) error {
	app := New()
	if err := app.Register(models...); err != nil {
		return err
	}
	return app.Start()
}
