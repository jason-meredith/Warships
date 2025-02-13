package main

import (
	"flag"
	"fmt"
	"github.com/jason-meredith/warships/game"
	"github.com/jason-meredith/warships/net"
	"os"
	"strconv"
	"strings"
	"time"
)

var logo = `
██╗    ██╗ █████╗ ██████╗ ███████╗██╗  ██╗██╗██████╗ ███████╗
██║    ██║██╔══██╗██╔══██╗██╔════╝██║  ██║██║██╔══██╗██╔════╝
██║ █╗ ██║███████║██████╔╝███████╗███████║██║██████╔╝███████╗
██║███╗██║██╔══██║██╔══██╗╚════██║██╔══██║██║██╔═══╝ ╚════██║
╚███╔███╔╝██║  ██║██║  ██║███████║██║  ██║██║██║     ███████║
 ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝     ╚══════╝
`

/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 22, 2018						 *
 *	FILE: 		launcher.go								 *
 *	PURPOSE:	Entry point to the main menu			 *
 *				 										 *
 *				The main function loads you into a 		 *
 *				menu giving you the option of Start 	 *
 *				Game, Join Game or Quit.				 *
 *														 *
 *				Selecting Start Game will launch an		 *
 *				options menu allowing you to set up		 *
 *				preferences for the new Game Server.	 *
 *														 *
 *				Join Game will prompt you to enter		 *
 *				server IP Address, Username and Password *
 *														 *
 *				Quit will exit the program				 *
 *														 *
 *														 *
 *********************************************************/

// main Launches the game into the main menu, showing logo, credits and initial options
func main() {

	// If the program is executed with args run in non-interactive start
	if len(os.Args) > 1 {
		nonInteractiveStartup(getArgs())

		// Otherwise run normally
	} else {

		clearScreen()
		fmt.Println(logo)
		fmt.Println("\t\t ~ By Jason Meredith ~ ")
		fmt.Println()

		inputMenu("Select an option\n1. Start Game\n2. Join Game\n3. Quit",
			func() { startServer() },
			func() { joinGame() },
			func() { os.Exit(0) })
	}

}

// getArgs gets all command line flags and returns them in a neatly packaged map
func getArgs() map[string]*string {

	var args map[string]*string
	args = make(map[string]*string)

	// Get flags

	args["mode"] = flag.String("mode", "client", "Starting a game or joining a game [ client | server ]")

	// Join flags
	args["joinAddress"] = flag.String("address", "127.0.0.1", "Server address to connect to")
	args["joinUsername"] = flag.String("username", "player", "Player username")
	args["password"] = flag.String("password", "", "Player password")

	args["hostAdminPassword"] = flag.String("admin-password", "", "Admin password")

	args["hostMaxPlayers"] = flag.String("max-players", "32", "Max players")
	args["hostShipLimit"] = flag.String("ship-limit", "16", "Ship limit")
	args["hostBoardSize"] = flag.String("board-size", "16", "Board size")
	args["deployPts"] = flag.String("deploy-points", "10", "Starting deployment points")

	commandMode := flag.Bool("cmd", false, "Run in single command mode")

	flag.Parse()

	msg := strings.Join(flag.Args(), " ")
	cmd := strconv.FormatBool(*commandMode)
	args["msg"] = &msg
	args["command"] = &cmd

	return args
}

// nonInteractiveStartup starts up warships based on the given command line flags
func nonInteractiveStartup(args map[string]*string) {

	// Run in server mode, starting a new game
	if *args["mode"] == "server" {

		maxPlayers, _ := strconv.Atoi(*args["hostMaxPlayers"])
		shipLimit, _ := strconv.Atoi(*args["hostShipLimit"])
		boardSize, _ := strconv.Atoi(*args["hostBoardSize"])
		deployPts, _ := strconv.Atoi(*args["deployPts"])

		newGame := game.Game{}
		newGame.Live = true
		newGame.Port = net.RPC_PORT
		newGame.Password = *args["password"]
		newGame.StartTime = time.Now()
		newGame.AdminPassword = *args["hostAdminPassword"]
		newGame.MaxPlayers = uint8(maxPlayers)
		newGame.ShipLimit = uint8(shipLimit)
		newGame.BoardSize = uint8(boardSize)
		newGame.Teams = []*game.Team{}
		newGame.StartDeployPts = deployPts

		net.StartGameServer(&newGame)

		// Run in client mode, connecting to an existing game
	} else if *args["mode"] == "client" {
		playerId, connection, err := net.CreateServerConnection(
			*args["joinUsername"],
			*args["password"],
			*args["joinAddress"])

		if err == nil {
			// If the -cmd flag is present, we take any remaining args after the officials ones
			// and run it like a player command. Otherwise simple join the game and continue
			// normally
			if *args["command"] == "true" {
				net.SendCommand(playerId, connection, *args["msg"])
			} else {
				net.AcceptCommands(playerId, connection)
			}
		} else {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
	}

}

// startServer shows the menu screen for starting a new server
func startServer() {

	const PASSWRD = "Password"
	const ADMIN_PASSWRD = "Admin Password"
	const MAX_PLAYERS = "Max Players"
	const SHIP_LIMIT = "Ship Limit"
	const BOARD_SIZE = "Board Size"
	const DEPLOY_POINTS = "Deployment Points"

	setupScreen()

	options := inputOptions(
		"Start New Game",
		MAX_PLAYERS,
		BOARD_SIZE,
		SHIP_LIMIT,
		PASSWRD,
		ADMIN_PASSWRD,
		DEPLOY_POINTS,
	)

	maxPlayers, err := strconv.Atoi(options[MAX_PLAYERS])
	shipLimit, err := strconv.Atoi(options[SHIP_LIMIT])
	boardSize, err := strconv.Atoi(options[BOARD_SIZE])
	deployPts, err := strconv.Atoi(options[DEPLOY_POINTS])

	if err != nil {
		// TODO: Handle this error pls
	}

	newGame := game.Game{}
	newGame.Live = true
	newGame.Port = net.RPC_PORT
	newGame.Password = options[PASSWRD]
	newGame.StartTime = time.Now()
	newGame.AdminPassword = options[ADMIN_PASSWRD]
	newGame.MaxPlayers = uint8(maxPlayers)
	newGame.ShipLimit = uint8(shipLimit)
	newGame.BoardSize = uint8(boardSize)
	newGame.Teams = []*game.Team{}
	newGame.StartDeployPts = deployPts

	clearScreen()
	net.StartGameServer(&newGame)

	// Proceed to game
}

// joinGame shows the menu screen for joining a running game server
func joinGame() {

	const SERV_ADDR = "Server Address"
	const PASSWRD = "Password"
	const USERNAME = "Username"

	setupScreen()

	success := false

	for !success {
		options := inputOptions(
			"Joining Game",
			SERV_ADDR,
			PASSWRD,
			USERNAME,
		)

		playerId, connection, err := net.CreateServerConnection(
			strings.TrimRight(options[USERNAME], "\n"),
			strings.TrimRight(options[PASSWRD], "\n"),
			options[SERV_ADDR])
		if err == nil {
			success = true
			net.AcceptCommands(playerId, connection)
		} else {
			fmt.Println("Error: " + err.Error())
			time.Sleep(700 * time.Millisecond)
			setupScreen()
		}
	}

}
