function init()
    width = 128
    height = 128
end

function update()
    --here u do ur game logic
end

function draw()
    --here u draw!
    --example draws a gradient
    for x = 0, width - 1 do
        for y = 0, height - 1 do
            local r = math.floor((x / (width - 1)) * 255)
            local g = 0
            local b = 0
            gfx.pset(x, y, r, g, b)
        end
    end
end