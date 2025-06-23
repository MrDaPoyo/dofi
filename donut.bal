let time = 0;
let screen_width = 128;
let screen_height = 128;
let A_speed = 0.05;
let B_speed = 0.03;
let R1 = 1.0;
let R2 = 2.0;
let K1 = 5.0;
let K2 = 5.0;
let cx = screen_width / 2;
let cy = screen_height / 2;

let _draw = fn() {
    let A = time * A_speed;
    let B = time * B_speed;
    let zbuf = array(screen_width * screen_height, 0.0);
    let buf = array(screen_width * screen_height, 0);

    let theta = 0.0;
    while (theta < 6.28) {
        let phi = 0.0;
        while (phi < 6.28) {
            let costheta = cos(theta);
            let sintheta = sin(theta);
            let cosphi = cos(phi);
            let sinphi = sin(phi);

            let circlex = R2 + R1 * costheta;
            let circley = R1 * sintheta;

            let cosB = cos(B);
            let sinB = sin(B);
            let x = circlex * cosB * cosphi - circley * sinB;
            let y = circlex * sinB * cosphi + circley * cosB;
            let z = K2 + circlex * sinphi;

            let ooz = 1.0 / z;
            let xp = (cx + K1 * ooz * x) as int;
            let yp = (cy + K1 * ooz * y) as int;

            let idx = xp + yp * screen_width;
            phi = phi + 0.07;
        }
        theta = theta + 0.07;
    }

    let y = 0;
    while (y < screen_height) {
        let x = 0;
        while (x < screen_width) {
            let idx = x + y * screen_width;
            let c = buf[idx];
            pset(x, y, c, c, c);
            x = x + 1;
        }
        y = y + 1;
    }

    time = time + 1;
    if (time > 128) {
        time = 0;
    }
}