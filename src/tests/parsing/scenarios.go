package parsing

var (
	ParsingScenarios = []struct {
		name    string
		args    string
		wantErr bool
	}{
		// Global test
		{
			name:    "Global - Invalid number of arguments (missing command)",
			args:    "my-identifier",
			wantErr: true,
		},

		// Create tests
		{
			name:    "Create - Valid topic",
			args:    "my-identifier create my-topic",
			wantErr: false,
		},
		{
			name:    "Create - Too many arguments",
			args:    "my-identifier create my-topic error",
			wantErr: true,
		},
		{
			name:    "Create - Missing topic argument",
			args:    "my-identifier create",
			wantErr: true,
		},
		{
			name:    "Create - Invalid option presence",
			args:    "my-identifier create my-topic --error",
			wantErr: true,
		},

		// List tests
		{
			name:    "List - Valid list",
			args:    "my-identifier list",
			wantErr: false,
		},
		{
			name:    "List - Too many arguments",
			args:    "my-identifier list error",
			wantErr: true,
		},
		{
			name:    "List - Invalid option presence",
			args:    "my-identifier list --error",
			wantErr: true,
		},

		// Produce tests
		{
			name:    "Produce - Valid topic",
			args:    "my-identifier produce my-topic",
			wantErr: false,
		},
		{
			name:    "Produce - Valid topic with option",
			args:    "my-identifier produce my-topic --raw",
			wantErr: false,
		},
		{
			name:    "Produce - Missing topic argument",
			args:    "my-identifier produce",
			wantErr: true,
		},
		{
			name:    "Produce - Too many arguments",
			args:    "my-identifier produce my-topic error",
			wantErr: true,
		},
		{
			name:    "Produce - Invalid option presence",
			args:    "my-identifier produce my-topic --error",
			wantErr: true,
		},

		// Subscribe tests
		{
			name:    "Subscribe - Valid topic and subscriber",
			args:    "my-identifier subscribe my-topic my-subscriber",
			wantErr: false,
		},
		{
			name:    "Subscribe - Valid with single option",
			args:    "my-identifier subscribe my-topic my-subscriber --raw",
			wantErr: false,
		},
		{
			name:    "Subscribe - Valid with multiple options",
			args:    "my-identifier subscribe my-topic my-subscriber --prefix my-prefix --raw",
			wantErr: false,
		},
		{
			name:    "Subscribe - Invalid option value placement",
			args:    "my-identifier subscribe my-topic my-subscriber --prefix --raw my-prefix",
			wantErr: true,
		},
		{
			name:    "Subscribe - Missing required arguments",
			args:    "my-identifier subscribe",
			wantErr: true,
		},
		{
			name:    "Subscribe - Too many arguments",
			args:    "my-identifier subscribe my-topic my-subscriber error",
			wantErr: true,
		},
		{
			name:    "Subscribe - Invalid option presence",
			args:    "my-identifier subscribe my-topic my-subscriber --error",
			wantErr: true,
		},

		// Info tests
		{
			name:    "Info - Get valid info",
			args:    "my-identifier info my-subscriber",
			wantErr: false,
		},
		{
			name:    "Info - Missing subscriber argument",
			args:    "my-identifier info",
			wantErr: true,
		},
		{
			name:    "Info - Too many arguments",
			args:    "my-identifier info my-subscriber error",
			wantErr: true,
		},
		{
			name:    "Info - Invalid option presence",
			args:    "my-identifier info my-subscriber --error",
			wantErr: true,
		},
	}
)
