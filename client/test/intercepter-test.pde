import oscP5.*;
import java.util.Date;

OscP5 oscP5;

void setup() {
  size(500, 500);
  oscP5 = new OscP5(this, 3333);//set port
}

/* incoming osc message are forwarded to the oscEvent method. */
void oscEvent(OscMessage theOscMessage) {
  /* print the address pattern and the typetag of the received OscMessage */
  println("### [" + new Date() + "]" + " received an osc message.");
  println(" addrpattern: "+theOscMessage.addrPattern());
  println(" typetag: "+theOscMessage.typetag());
  println(" get(0): "+theOscMessage.get(0).intValue());
}

void draw() {
  background(255);
  noStroke();
  fill(0, 0, 0);
}
