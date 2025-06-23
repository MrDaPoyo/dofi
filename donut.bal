let _draw = fn() {
    let A = 0.0;
    let B = 0.0;
    let width = 80;
    let height = 40;
    let theta = 0.0;
    let phi = 0.0;
    let two_pi = 6.28318;
    let increment_theta = 0.07;
    let increment_phi = 0.02;
    let R1 = 1.0;
    let R2 = 2.0;
    let K2 = 5.0;
    let K1 = width * K2 * 3.0 / (8.0 * (R1 + R2));

    phi = 0.0;
    while (phi < two_pi) {
        theta = 0.0;
        while (theta < two_pi) {
            let costheta = cos(theta);
            let sintheta = sin(theta);
            let cosphi = cos(phi);
            let sinphi = sin(phi);

            let circlex = R2 + R1 * costheta;
            let circley = R1 * sintheta;

            let x = circlex * (cosB = cos(B)) * cosphi - circley * sin(B);
            let y = circlex * sin(B) * cosphi + circley * cosB;
            let z = K2 + circlex * sinphi;

            let ooz = 1.0 / z;
            let xp = (width / 2.0 + K1 * ooz * x) as int;
            let yp = (height / 2.0 - K1 * ooz * y) as int;

            let L = cosphi * costheta * sin(B) - cos(B) * costheta * sinphi - sintheta * cos(B);
            let brightness = (L > 0.0) ? (L * 255.0) as int : 0;

            if (xp >= 0) {
                if (xp < width) {
                    if (yp >= 0) {
                        if (yp < height) {
                            pset(xp, yp, brightness, brightness, brightness);
                        }
                    }
                }
            }

            theta = theta + increment_theta;
        }
        phi = phi + increment_phi;
    }
};

let _update = fn() {};