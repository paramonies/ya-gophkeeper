package cmd

type Command string

const (
	RootCommand Command = "ROOT"

	RegisterUserCommand Command = "REGISTER_USER"
	LoginUserCommand    Command = "LOGIN_USER"
	SyncDataCommand     Command = "SYNC_DATA"

	CreatePasswordCommand Command = "CREATE_PASSWORD"
	GetPasswordCommand    Command = "GET_PASSWORD"
	DeletePasswordCommand Command = "DELETE_PASSWORD"

	CreateTextCommand Command = "CREATE_TEXT"
	GetTextCommand    Command = "GET_TEXT"
	DeleteTextCommand Command = "DELETE_TEXT"

	CreateBinaryCommand Command = "CREATE_BINARY"
	GetBinaryCommand    Command = "GET_BINARY"
	DeleteBinaryCommand Command = "DELETE_BINARY"
)

var commands = map[Command]struct {
	Use   string
	Short string
	Long  string
}{
	RootCommand: {
		Use:   "gophkeeper",
		Short: "GophKeeper is a service to store and protect card, binary, password, text data",
		Long: "GophKeeper is a service, that gives you possibilities to save card, binary, password, text data on server and and get up-to-date data on various clients.\n" +
			"Service is synchronized between all you devices, where you are authenticated. This application is a CLI tool to interact with the service.\nType --help to see more.",
	},

	RegisterUserCommand: {
		Use:   "registerUser",
		Short: "Register new user in service",
		Long:  "This command register a new user.\nUsage: gophkeeperclient registerUser --login=<login> --password=<password>",
	},
	LoginUserCommand: {
		Use:   "loginUser",
		Short: "Login registered user in service",
		Long:  "This command login user.\nUsage: gophkeeperclient loginUser --login=<login> --password=<password>",
	},
	SyncDataCommand: {
		Use:   "syncUserData",
		Short: "Synchronize local user data with the server database",
		Long: "This command provides latest data from the server.\n" +
			"if server data has version higher that the local ones it saves in local storage.\n" +
			"During the saving of local data to the server in case of version conflict(server data version is higher/newer)" +
			"you will be alerted by a warning.\n" +
			"Usage: gophkeeperclient syncUserData",
	},

	CreatePasswordCommand: {
		Use:   "createPassword",
		Short: "Create a new password entity for uniq login",
		Long: "This command allows authenticated user to save new password data.\n" +
			"Usage: gophkeeperclient createPassword --login=<login> --password=<password> --meta=<meta_info>",
	},
	GetPasswordCommand: {
		Use:   "getPassword",
		Short: "Get password entity for login",
		Long: "This command returns password data for requested login and authenticated user.\n" +
			"Usage: gophkeeperclient getPassword --login=<login>",
	},
	DeletePasswordCommand: {
		Use:   "deletePassword",
		Short: "Delete password entity for login",
		Long: "This command allows authenticated user to delete the password data for login.\n" +
			"Usage: gophkeeperclient deletePassword --login=<login>",
	},

	CreateTextCommand: {
		Use:   "createText",
		Short: "Create a new text entity for uniq title",
		Long: "This command allows authenticated user to save new text data.\n" +
			"Usage: gophkeeperclient createText --title=<title> --data=<data> --meta=<meta_info>",
	},
	GetTextCommand: {
		Use:   "getText",
		Short: "Get text entity for title",
		Long: "This command returns text data for requested title and authenticated user.\n" +
			"Usage: gophkeeperclient getText --title=<title>",
	},
	DeleteTextCommand: {
		Use:   "deleteText",
		Short: "Delete text entity for title",
		Long: "This command allows authenticated user to delete the text data for title.\n" +
			"Usage: gophkeeperclient deleteText --title=<title>",
	},

	CreateBinaryCommand: {
		Use:   "createBinary",
		Short: "Create a new binary entity for uniq title",
		Long: "This command allows authenticated user to save new binary data.\n" +
			"Usage: gophkeeperclient createBinary --title=<title> --data=<data> --meta=<meta_info>",
	},
	GetBinaryCommand: {
		Use:   "getBinary",
		Short: "Get binary entity for title",
		Long: "This command returns binary data for requested title and authenticated user.\n" +
			"Usage: gophkeeperclient getBinary --title=<title>",
	},
	DeleteBinaryCommand: {
		Use:   "deleteBinary",
		Short: "Delete binary entity for title",
		Long: "This command allows authenticated user to delete the binary data for title.\n" +
			"Usage: gophkeeperclient deleteBinary --title=<title>",
	},
}
