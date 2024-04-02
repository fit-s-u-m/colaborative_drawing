let socket = new WebSocket("ws://localhost:8080/ws");
console.log("attempting to connect");

socket.onclose = () => {
  console.log("disconnected");
};
socket.onerror = (err) => {
  console.log("error", err);
};
socket.onopen = () => {
  console.log("connected");
};


socket.onmessage = (mess) => {   // receive messages
  const data = JSON.parse(mess.data)
    // if drawing have been send 
    fill("blue");
    noStroke();
    circle(data.x, data.y, 30);

   // create a new div put it right to show you are receiving it
    if(data.chat){ // if chat is send
      const father = createDiv().parent("chat-box")
                     .class("chat chat-end")
        createDiv(data.chat).parent(father)
                       .class("chat-bubble chat-bubble-accent mt-3")
    }
};

function setup() {
  createCanvas(windowWidth/2, windowHeight/2).parent("canvas");
  background(0);
  const button = createButton('Clear');
  button.mousePressed(()=>{background(0)})
         .parent("clear")
}
function draw() {
}


// ----------- for drawing
function mouseDragged() { 
  fill(255);
  noStroke();
  circle(mouseX, mouseY, 30);
  socket.send(JSON.stringify({ x: mouseX, y: mouseY })); // end to the back-end to broadcast the signal
}


// ----------- for chat
function send(event){
  event.preventDefault();
  const chat = select("#chat").value()

  var input = document.getElementById("chat");
  input.value ="";

  socket.send(JSON.stringify({chat:chat}));

  // create a new div put it right to show you are sending it
  const father = createDiv().parent("chat-box")
                 .class("chat chat-start m-0")
  createDiv(chat).parent(father)
                 .class("chat-bubble chat-bubble-primary mt-3")
}
