var time = 0;
var screen_width = 128;
var screen_height = 128;
var cx = screen_width / 2;
var cy = screen_height / 2;

var R1 = 20;  // Minor radius (tube thickness)
var R2 = 40;  // Major radius (donut size)
var K2 = 200; // Distance from viewer
var K1 = screen_width * K2 * 3 / (8 * (R1 + R2));

// Rotation speeds
var A_speed = 0.07;
var B_speed = 0.03;

fn clamp(v, lo, hi) {
  if (v < lo) {
    return lo;
  }
  if (v > hi) {
    return hi;
  }
  return v;
}

fn _update() {
  time = time + 1;
}

fn _draw() {
  cls();

  var A = time * A_speed;
  var B = time * B_speed;

  var cosA = cos(A);
  var sinA = sin(A);
  var cosB = cos(B);
  var sinB = sin(B);

  var theta_spacing = 0.07;
  var phi_spacing = 0.02;

  var theta = 0.0;
  while (theta < 6.28) {
    var costheta = cos(theta);
    var sintheta = sin(theta);

    var phi = 0.0;
    while (phi < 6.28) {
      var cosphi = cos(phi);
      var sinphi = sin(phi);
      var circlex = R2 + R1 * costheta;
      var circley = R1 * sintheta;
      var x = circlex * (cosB * cosphi + sinA * sinB * sinphi) - circley * cosA * sinB;
      var y = circlex * (sinB * cosphi - sinA * cosB * sinphi) + circley * cosA * cosB;
      var z = K2 + cosA * circlex * sinphi + circley * sinA;
      var ooz = 1 / z; // "one over z"
      var xp = floor(cx + K1 * ooz * x + 0.5);
      var yp = floor(cy - K1 * ooz * y + 0.5);
      var L = cosphi * costheta * sinB - cosA * costheta * sinphi - sinA * sintheta + cosB * (cosA * sintheta - costheta * sinA * sinphi);
      if (L > 0 && xp >= 0 && xp < screen_width && yp >= 0 && yp < screen_height) {
        var clamped_L = clamp(L, 0.1, 1);
        var r = floor(255 * clamped_L);
        var g = floor(255 * clamped_L * 0.8);
        var b = floor(255 * clamped_L * 0.6);
        pset(xp, yp, r, g, b);
      }
      phi = phi + phi_spacing;
    }
    theta = theta + theta_spacing;
  }
}

fn _init() {
  print("3D Spinning Donut initialized!");
}

_init();
