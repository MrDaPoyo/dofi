var A = 0;
var B = 0;

fn _draw() {
  var cosA = cos(A);
  var sinA = sin(A);
  var cosB = cos(B);
  var sinB = sin(B);
  
  pset(64, 64, 255, 0, 0);
  pset(65, 64, 0, 255, 0);
  pset(64, 65, 0, 0, 255);
  
  pset(84, 64, 255, 255, 255);  // right
  pset(78, 78, 255, 255, 255);  // bottom-right
  pset(64, 84, 255, 255, 255);  // bottom
  pset(50, 78, 255, 255, 255);  // bottom-left
  pset(44, 64, 255, 255, 255);  // left
  pset(50, 50, 255, 255, 255);  // top-left
  pset(64, 44, 255, 255, 255);  // top
  pset(78, 50, 255, 255, 255);  // top-right
  
  var testX = 0;
  while (testX < 5) {
    pset(testX, 10, 255, 0, 255);  // magenta pixels
    testX = testX + 1;
  }
  
  while (testY < 10) {
    pset(testY, 20, 0, 255, 255);  // cyan pixels
    testY = testY + 1
  }
  
  A = A + 0.04;
  B = B + 0.02;
}