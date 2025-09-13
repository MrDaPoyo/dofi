function init()
    width = 128
    height = 128
end

function update()
end

function draw()
    for x = 0, width - 1 do
        for y = 0, height - 1 do
            local value = ((x + y) % 2) * 255
            gfx.pset(x, y, value, value, value)
        end
    end
end
