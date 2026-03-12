import android.accessibilityservice.AccessibilityServiceInfo;
import android.app.UiAutomation;
import android.os.HandlerThread;
import android.os.Looper;
import android.view.accessibility.AccessibilityNodeInfo;

import java.io.PrintStream;
import java.lang.reflect.Constructor;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import java.util.List;

/**
 * Lightweight helper that runs via app_process on Android device.
 * Connects to accessibility framework and reads UI text
 * from live stream apps, skipping video surfaces to avoid dump timeout.
 *
 * Usage:
 *   adb push classes.dex /data/local/tmp/adbclaw-monitor.dex
 *   adb shell CLASSPATH=/data/local/tmp/adbclaw-monitor.dex app_process / ADBClawMonitor [--interval 2000] [--count 0]
 */
public class ADBClawMonitor {

    private static final PrintStream out = System.out;
    private static final PrintStream err = System.err;

    // Classes known to be video surfaces — skip traversing into these
    private static final String[] SKIP_CLASSES = {
        "TextureRenderView",
        "TextureView",
        "SurfaceView",
        "VideoView",
        "GLSurfaceView",
    };

    public static void main(String[] args) {
        int interval = 2000;
        int maxCount = 0; // 0 = unlimited

        for (int i = 0; i < args.length; i++) {
            if ("--interval".equals(args[i]) && i + 1 < args.length) {
                interval = Integer.parseInt(args[++i]);
            } else if ("--count".equals(args[i]) && i + 1 < args.length) {
                maxCount = Integer.parseInt(args[++i]);
            }
        }

        err.println("[ADBClawMonitor] Starting... interval=" + interval + "ms, count=" + maxCount);

        UiAutomation uiAutomation = null;
        try {
            uiAutomation = connectUiAutomation();
        } catch (Exception e) {
            err.println("[ADBClawMonitor] Failed to connect: " + e.getMessage());
            e.printStackTrace(err);
            System.exit(1);
        }

        err.println("[ADBClawMonitor] Connected to accessibility framework");

        // Wait for accessibility framework to fully initialize
        sleep(1000);

        final UiAutomation uiAutoFinal = uiAutomation;
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            try {
                disconnectUiAutomation(uiAutoFinal);
                err.println("[ADBClawMonitor] Disconnected");
            } catch (Exception ignored) {}
        }));

        LinkedHashSet<String> seenMessages = new LinkedHashSet<>();
        int pollCount = 0;

        try {
            while (maxCount == 0 || pollCount < maxCount) {
                pollCount++;
                try {
                    AccessibilityNodeInfo root = uiAutomation.getRootInActiveWindow();
                    if (root == null) {
                        err.println("[ADBClawMonitor] Root node is null, retrying...");
                        sleep(interval);
                        continue;
                    }

                    List<String[]> messages = new ArrayList<>();
                    findTextNodes(root, messages);
                    root.recycle();

                    for (String[] msg : messages) {
                        String key = msg[0] + "|" + msg[1];
                        if (!seenMessages.contains(key)) {
                            seenMessages.add(key);
                            out.println("{\"text\":" + jsonString(msg[0])
                                + ",\"class\":" + jsonString(msg[1]) + "}");
                            out.flush();
                        }
                    }

                    // Prevent unbounded growth
                    if (seenMessages.size() > 500) {
                        List<String> list = new ArrayList<>(seenMessages);
                        seenMessages.clear();
                        seenMessages.addAll(list.subList(list.size() - 200, list.size()));
                    }

                } catch (Exception e) {
                    err.println("[ADBClawMonitor] Error: " + e.getMessage());
                }

                sleep(interval);
            }
        } finally {
            try { disconnectUiAutomation(uiAutomation); } catch (Exception ignored) {}
            err.println("[ADBClawMonitor] Done after " + pollCount + " polls");
        }
    }

    /**
     * Connect to UiAutomation using reflection for hidden APIs.
     * This is the same mechanism uiautomator itself uses.
     */
    @SuppressWarnings("unchecked")
    private static UiAutomation connectUiAutomation() throws Exception {
        HandlerThread ht = new HandlerThread("UiAutoThread");
        ht.start();

        // android.app.UiAutomationConnection is @hide, extends IUiAutomationConnection$Stub
        Class<?> connClass = Class.forName("android.app.UiAutomationConnection");
        Object conn = connClass.getDeclaredConstructor().newInstance();

        // Constructor is UiAutomation(Looper, IUiAutomationConnection)
        Class<?> iConnClass = Class.forName("android.app.IUiAutomationConnection");
        Constructor<UiAutomation> ctor = UiAutomation.class.getDeclaredConstructor(
            Looper.class, iConnClass
        );
        ctor.setAccessible(true);
        UiAutomation uiAuto = ctor.newInstance(ht.getLooper(), conn);

        // UiAutomation.connect() is @hide
        Method connectMethod = UiAutomation.class.getDeclaredMethod("connect");
        connectMethod.setAccessible(true);
        connectMethod.invoke(uiAuto);

        // Enable multi-window access
        try {
            AccessibilityServiceInfo info = uiAuto.getServiceInfo();
            if (info != null) {
                info.flags |= AccessibilityServiceInfo.FLAG_RETRIEVE_INTERACTIVE_WINDOWS;
                uiAuto.setServiceInfo(info);
            }
        } catch (Exception e) {
            err.println("[ADBClawMonitor] Warning: set service info failed: " + e.getMessage());
        }

        return uiAuto;
    }

    /**
     * Recursively find text-bearing nodes, skipping video surface subtrees.
     */
    private static void findTextNodes(AccessibilityNodeInfo node, List<String[]> results) {
        if (node == null) return;

        String className = node.getClassName() != null ? node.getClassName().toString() : "";

        // Skip video surface subtrees entirely
        for (String skip : SKIP_CLASSES) {
            if (className.contains(skip)) {
                return;
            }
        }

        CharSequence text = node.getText();
        CharSequence desc = node.getContentDescription();

        if (text != null && text.length() > 0) {
            results.add(new String[]{text.toString(), className});
        } else if (desc != null && desc.length() > 0) {
            results.add(new String[]{desc.toString(), className});
        }

        int childCount = node.getChildCount();
        for (int i = 0; i < childCount; i++) {
            AccessibilityNodeInfo child = node.getChild(i);
            if (child != null) {
                findTextNodes(child, results);
                child.recycle();
            }
        }
    }

    private static void disconnectUiAutomation(UiAutomation uiAuto) {
        try {
            Method m = UiAutomation.class.getDeclaredMethod("disconnect");
            m.setAccessible(true);
            m.invoke(uiAuto);
        } catch (Exception ignored) {}
    }

    private static String jsonString(String s) {
        if (s == null) return "null";
        return "\"" + s.replace("\\", "\\\\")
                        .replace("\"", "\\\"")
                        .replace("\n", "\\n")
                        .replace("\r", "\\r")
                        .replace("\t", "\\t")
               + "\"";
    }

    private static void sleep(int ms) {
        try { Thread.sleep(ms); } catch (InterruptedException ignored) {}
    }
}
