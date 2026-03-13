import android.media.AudioFormat;
import android.media.AudioRecord;

import java.io.OutputStream;
import java.io.PrintStream;

/**
 * Lightweight helper that runs via app_process on Android device.
 * Captures system audio via REMOTE_SUBMIX and outputs raw PCM to stdout.
 *
 * Requires Android 11+ (API 30). REMOTE_SUBMIX is accessible to shell user
 * (uid 2000) when running through app_process.
 *
 * WARNING: REMOTE_SUBMIX mutes device speakers while capturing.
 *
 * Usage:
 *   adb push classes.dex /data/local/tmp/adbclaw-audio.dex
 *   adb exec-out CLASSPATH=/data/local/tmp/adbclaw-audio.dex \
 *       app_process / ADBClawAudio [--rate 16000] [--duration 0]
 *
 * Output: WAV header (44 bytes) + raw PCM 16-bit mono to stdout.
 * Duration 0 = unlimited (until killed or stdin closed).
 */
public class ADBClawAudio {

    private static final PrintStream err = System.err;

    // MediaRecorder.AudioSource.REMOTE_SUBMIX = 8
    private static final int AUDIO_SOURCE_REMOTE_SUBMIX = 8;

    private static volatile boolean running = true;

    public static void main(String[] args) {
        int sampleRate = 16000;
        int durationMs = 0; // 0 = unlimited

        for (int i = 0; i < args.length; i++) {
            if ("--rate".equals(args[i]) && i + 1 < args.length) {
                sampleRate = Integer.parseInt(args[++i]);
            } else if ("--duration".equals(args[i]) && i + 1 < args.length) {
                durationMs = Integer.parseInt(args[++i]);
            }
        }

        err.println("[ADBClawAudio] Starting... rate=" + sampleRate + "Hz, duration=" + durationMs + "ms");
        err.println("[ADBClawAudio] Using REMOTE_SUBMIX (device speakers will be muted)");

        // Check Android version
        int sdkInt = android.os.Build.VERSION.SDK_INT;
        err.println("[ADBClawAudio] Android SDK: " + sdkInt);
        if (sdkInt < 30) {
            err.println("[ADBClawAudio] ERROR: Requires Android 11+ (API 30), current: " + sdkInt);
            System.exit(1);
        }

        int channelConfig = AudioFormat.CHANNEL_IN_MONO;
        int audioFormat = AudioFormat.ENCODING_PCM_16BIT;
        int bufferSize = AudioRecord.getMinBufferSize(sampleRate, channelConfig, audioFormat);
        if (bufferSize <= 0) {
            err.println("[ADBClawAudio] ERROR: getMinBufferSize returned " + bufferSize);
            System.exit(1);
        }
        // Use at least 4x minimum buffer for smoother streaming
        bufferSize = Math.max(bufferSize, sampleRate * 2 /* 16-bit = 2 bytes */ * 1 /* 1 second */);

        AudioRecord recorder;
        try {
            recorder = new AudioRecord(
                AUDIO_SOURCE_REMOTE_SUBMIX,
                sampleRate,
                channelConfig,
                audioFormat,
                bufferSize
            );
        } catch (Exception e) {
            err.println("[ADBClawAudio] ERROR: Failed to create AudioRecord: " + e.getMessage());
            err.println("[ADBClawAudio] This may indicate REMOTE_SUBMIX is not available on this device");
            System.exit(1);
            return;
        }

        if (recorder.getState() != AudioRecord.STATE_INITIALIZED) {
            err.println("[ADBClawAudio] ERROR: AudioRecord failed to initialize (state=" + recorder.getState() + ")");
            err.println("[ADBClawAudio] REMOTE_SUBMIX may not be available or permitted");
            recorder.release();
            System.exit(1);
        }

        // Graceful shutdown on SIGTERM / process kill
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            running = false;
            err.println("[ADBClawAudio] Shutting down...");
        }));

        try {
            recorder.startRecording();
        } catch (Exception e) {
            err.println("[ADBClawAudio] ERROR: startRecording failed: " + e.getMessage());
            recorder.release();
            System.exit(1);
            return;
        }

        err.println("[ADBClawAudio] Recording started");

        OutputStream out = System.out;
        byte[] readBuffer = new byte[bufferSize];

        try {
            // Write WAV header (we'll use a streaming-friendly approach:
            // set data size to max value since we don't know duration upfront)
            writeWavHeader(out, sampleRate);

            long startTime = System.currentTimeMillis();
            long totalBytesWritten = 0;

            while (running) {
                // Check duration limit
                if (durationMs > 0 && (System.currentTimeMillis() - startTime) >= durationMs) {
                    err.println("[ADBClawAudio] Duration limit reached");
                    break;
                }

                int bytesRead = recorder.read(readBuffer, 0, readBuffer.length);
                if (bytesRead > 0) {
                    out.write(readBuffer, 0, bytesRead);
                    out.flush();
                    totalBytesWritten += bytesRead;
                } else if (bytesRead == AudioRecord.ERROR_INVALID_OPERATION) {
                    err.println("[ADBClawAudio] ERROR: Invalid operation during read");
                    break;
                } else if (bytesRead == AudioRecord.ERROR_BAD_VALUE) {
                    err.println("[ADBClawAudio] ERROR: Bad value during read");
                    break;
                } else if (bytesRead == AudioRecord.ERROR) {
                    err.println("[ADBClawAudio] ERROR: Generic error during read");
                    break;
                }
                // bytesRead == 0 is possible under pressure, just continue
            }

            long elapsed = System.currentTimeMillis() - startTime;
            err.println("[ADBClawAudio] Done: " + totalBytesWritten + " bytes, " + elapsed + "ms");
        } catch (Exception e) {
            // Broken pipe is normal when the host process closes the connection
            if (!e.getMessage().contains("Broken pipe")) {
                err.println("[ADBClawAudio] ERROR: " + e.getMessage());
            }
        } finally {
            try { recorder.stop(); } catch (Exception ignored) {}
            recorder.release();
            err.println("[ADBClawAudio] Released AudioRecord");
        }
    }

    /**
     * Write a WAV header for streaming PCM data.
     * Uses 0x7FFFFFFF for data size since total length is unknown.
     */
    private static void writeWavHeader(OutputStream out, int sampleRate) throws Exception {
        int channels = 1;
        int bitsPerSample = 16;
        int byteRate = sampleRate * channels * bitsPerSample / 8;
        int blockAlign = channels * bitsPerSample / 8;
        // Use max int32 for unknown-length streaming
        int dataSize = 0x7FFFFFFF;
        int chunkSize = 36 + dataSize;

        byte[] header = new byte[44];
        // RIFF header
        header[0] = 'R'; header[1] = 'I'; header[2] = 'F'; header[3] = 'F';
        writeInt32LE(header, 4, chunkSize);
        header[8] = 'W'; header[9] = 'A'; header[10] = 'V'; header[11] = 'E';
        // fmt subchunk
        header[12] = 'f'; header[13] = 'm'; header[14] = 't'; header[15] = ' ';
        writeInt32LE(header, 16, 16); // subchunk1 size (PCM = 16)
        writeInt16LE(header, 20, (short) 1); // audio format (PCM = 1)
        writeInt16LE(header, 22, (short) channels);
        writeInt32LE(header, 24, sampleRate);
        writeInt32LE(header, 28, byteRate);
        writeInt16LE(header, 32, (short) blockAlign);
        writeInt16LE(header, 34, (short) bitsPerSample);
        // data subchunk
        header[36] = 'd'; header[37] = 'a'; header[38] = 't'; header[39] = 'a';
        writeInt32LE(header, 40, dataSize);

        out.write(header);
        out.flush();
    }

    private static void writeInt32LE(byte[] buf, int offset, int value) {
        buf[offset]     = (byte) (value & 0xFF);
        buf[offset + 1] = (byte) ((value >> 8) & 0xFF);
        buf[offset + 2] = (byte) ((value >> 16) & 0xFF);
        buf[offset + 3] = (byte) ((value >> 24) & 0xFF);
    }

    private static void writeInt16LE(byte[] buf, int offset, short value) {
        buf[offset]     = (byte) (value & 0xFF);
        buf[offset + 1] = (byte) ((value >> 8) & 0xFF);
    }
}
