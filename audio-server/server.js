const http = require('http');
const child_process = require('child_process');
const url = require('url');

const server = http.createServer((req, res) => {
    const requestUrl = url.parse(req.url, true);

    if (requestUrl.pathname === '/stream') {
        const ffmpeg = child_process.spawn('ffmpeg', ['-ar', '44100', '-ac', '1', '-f', 'alsa', '-i', 'plughw:1,0', '-f', 'wav', 'pipe:1']);
        res.setHeader('Content-Type', 'audio/wav');

        ffmpeg.stdout.pipe(res);

        // Handle client disconnect
        req.on('close', () => {
            ffmpeg.kill('SIGTERM');
        });

        // Handle ffmpeg error
        ffmpeg.on('error', (err) => {
            console.error('ffmpeg error:', err);
            ffmpeg.kill('SIGTERM');
        });

        // Handle ffmpeg exit
        ffmpeg.on('exit', (code, signal) => {
            if (code !== null) {
                console.error(`ffmpeg exited with code ${code}`);
            }
            if (signal !== null) {
                console.error(`ffmpeg was killed with signal ${signal}`);
            }
        });
    } else {
        res.statusCode = 404;
        res.end('Not Found');
    }
});

server.listen(8081, () => {
    console.log('Server is listening on port 8081');
});

// Handle server error
server.on('error', (err) => {
    console.error('Server error:', err);
});
