-- 3D Spinning Donut
local time = 0
local screen_width = 128
local screen_height = 128
local cx = screen_width / 2
local cy = screen_height / 2

-- Donut parameters
local R1 = 20  -- Minor radius (tube thickness)
local R2 = 40  -- Major radius (donut size)
local K2 = 200 -- Distance from viewer
local K1 = screen_width * K2 * 3 / (8 * (R1 + R2))

-- Rotation speeds
local A_speed = 0.07
local B_speed = 0.03

function _update()
    time = time + 1
end

function _draw()
    dofi.cls()
    
    local A = time * A_speed
    local B = time * B_speed
    
    -- Precompute trig values
    local cosA = math.cos(A)
    local sinA = math.sin(A)
    local cosB = math.cos(B)
    local sinB = math.sin(B)
    
    local zbuffer = {}
    local output = {}
    
    for i = 0, screen_width * screen_height - 1 do
        zbuffer[i] = 0
        output[i] = 32
    end
    
    local theta_spacing = 0.07
    local phi_spacing = 0.02
    
    local theta = 0
    while theta < 6.28 do
        local costheta = math.cos(theta)
        local sintheta = math.sin(theta)
        
        local phi = 0
        while phi < 6.28 do
            local cosphi = math.cos(phi)
            local sinphi = math.sin(phi)
            
            -- 3D coordinates of the torus
            local circlex = R2 + R1 * costheta
            local circley = R1 * sintheta
            
            -- 3D coordinates after rotation
            local x = circlex * (cosB * cosphi + sinA * sinB * sinphi) - circley * cosA * sinB
            local y = circlex * (sinB * cosphi - sinA * cosB * sinphi) + circley * cosA * cosB
            local z = K2 + cosA * circlex * sinphi + circley * sinA
            local ooz = 1 / z -- "one over z"
            
            -- Project to 2D with better rounding
            local xp = math.floor(cx + K1 * ooz * x + 0.5)
            local yp = math.floor(cy - K1 * ooz * y + 0.5)
            
            -- Calculate luminance
            local L = cosphi * costheta * sinB - cosA * costheta * sinphi - sinA * sintheta + cosB * (cosA * sintheta - costheta * sinA * sinphi)
            
            if L > 0 and xp >= 0 and xp < screen_width and yp >= 0 and yp < screen_height then
                local idx = yp * screen_width + xp
                if ooz > zbuffer[idx] then
                    zbuffer[idx] = ooz
                    
                    -- Color based on luminance with proper clamping and minimum brightness
                    local clamped_L = math.max(0.1, math.min(1, L))
                    local r = math.floor(255 * clamped_L)
                    local g = math.floor(255 * clamped_L * 0.8)
                    local b = math.floor(255 * clamped_L * 0.6)
                    
                    dofi.pset(xp, yp, r, g, b)
                end
            end
            
            phi = phi + phi_spacing
        end
        theta = theta + theta_spacing
    end
end

function _init()
    print("3D Spinning Donut initialized!")
end

if _init then
    _init()
end