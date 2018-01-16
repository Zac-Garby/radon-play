var sock

// opens a new websocket and sends the code
// to the server, waiting for response(s)
function run() {
	clear()

	if (sock !== undefined) {
		sock.close()
	}
	
	sock = new WebSocket(`ws://${window.location.host}/run`)
	write("waiting for server response...\n")

	sock.onmessage = handleMessage
	sock.onopen = sendCode
}

// handles a message from the server.
// write()s it to the output
function handleMessage(evt) {
	write(evt.data)
}

// sends the value of #input to the server
function sendCode(evt) {
	var textarea = document.getElementById("input")
	var text = textarea.value

	sock.send(text)
}

// writes some text to the output
function write(text) {
	var out = document.getElementById("out")
	out.value += text
}

// clears the output
function clear() {
	var out = document.getElementById("out")
	out.value = ""
}

// sets up the site -- sets the default values of #input and #out
function setup() {
	document.getElementById("input").value = `fact(n) = match n where
    | 0 -> 1,
    | _ -> n * fact(n - 1)

for i = 0; i <= 15; i = i + 1 do
    print(fact(i))
`

	document.getElementById("out").value = "press ► to run your code"
}
