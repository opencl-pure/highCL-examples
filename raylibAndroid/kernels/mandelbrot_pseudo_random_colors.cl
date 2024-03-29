__kernel void mandelbrot_pseudo_random_colors(__write_only image2d_t outputImage)
{
    const float zoom = 0.004;
    const float moveX = -1.8;
    const float moveY = -1.1;
    const int max_iteration = 200;

    int2 pos = (int2)(get_global_id(0), get_global_id(1));
    int width = get_global_size(0);
    int height = get_global_size(1);

    float2 c = (float2)(pos.x * zoom + moveX, pos.y * zoom + moveY);

    float2 z = 0.0f;
    float2 zn = 0.0f;

    int iteration = 0;
    float random_value = 0.0f;

    while(iteration < max_iteration && length(z) < 2.0f)
    {
        zn.x = z.x * z.x - z.y * z.y + c.x;
        zn.y = 2 * z.x * z.y + c.y;

        z = zn;

        iteration++;

        // Perform modulo operation to get a value between 0 and 1
        random_value = ((float)((int)iteration*1000 % (int)max_iteration))/1000.0;

        // Normalize the value to the range [0, 1]
        random_value /= (float)max_iteration;
    }

    float ratio = (float)iteration / (float)max_iteration;
    float4 color;

    if (iteration == max_iteration) {
        color = (float4)(0.0f, 0.0f, 0.0f, 1.0f);
    }
    else {
        float value = log((float)iteration) / log((float)max_iteration);
        float angle = value * 3.1415927 * 2.0 + 3.1415927 * 0.5;
        color.x = cos(angle) * value;
        color.y = sin(angle) * value;
        color.z = random_value;
        color.w = 1.0f;
    }

    write_imagef(outputImage, pos, color);
}
