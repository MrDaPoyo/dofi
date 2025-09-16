-- welcome to dofi!
-------------------
-- this is an example script.
-- try executing it! (esc to exit)

function init()
    width = 128
    height = 128
    rotation = 0
end

function update()
    rotation = rotation + 0.05
end

function draw()
    --here u draw!
    --example draws a gradient
    local cx = width / 2
    local cy = height / 2

    local cosr = math.cos(rotation)
    local sinr = math.sin(rotation)

    for x = 0, width - 1 do
        for y = 0, height - 1 do
            -- translate pixel to center
            local dx = x - cx
            local dy = y - cy

            -- rotate coordinates
            local rx = dx * cosr - dy * sinr
            local ry = dx * sinr + dy * cosr

            -- map rotated x to color gradient
            local r = math.floor(((rx + cx) / width) * 255)
            local g = math.floor(((ry + cy) / height) * 255)
            local b = 128

            -- clamp to [0, 255]
            if r < 0 then r = 0 elseif r > 255 then r = 255 end
            if g < 0 then g = 0 elseif g > 255 then g = 255 end

            gfx.pset(x, y, r, g, b)
        end
    end
end