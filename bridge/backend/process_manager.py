import subprocess
import os
import signal
import threading
import json
import time
import re

class ProcessManager:
    def __init__(self, executable_path="./fortune"):
        self.executable_path = executable_path
        self.process = None
        self.status = {"running": False, "last_heartbeat": "", "found_keys": [], "total_tried": 0, "iops": 0}
        self.lock = threading.Lock()
        self.log_file = "scanner.log"

    def start(self, args):
        with self.lock:
            if self.process and self.process.poll() is None:
                return False, "Scanner is already running"

            cmd = [self.executable_path] + args

            try:
                self.process = subprocess.Popen(
                    cmd,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.STDOUT,
                    text=True,
                    bufsize=1,
                    universal_newlines=True,
                    preexec_fn=os.setsid
                )
                self.status["running"] = True
                self.status["total_tried"] = 0
                self.status["iops"] = 0
                threading.Thread(target=self._read_output, daemon=True).start()
                return True, "Scanner started"
            except Exception as e:
                return False, str(e)

    def stop(self):
        with self.lock:
            if not self.process or self.process.poll() is not None:
                self.status["running"] = False
                return False, "Scanner is not running"

            try:
                os.killpg(os.getpgid(self.process.pid), signal.SIGTERM)
                self.process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                os.killpg(os.getpgid(self.process.pid), signal.SIGKILL)

            self.status["running"] = False
            return True, "Scanner stopped"

    def _read_output(self):
        # Regex for parsing heartbeat output like "tried 24468 keys (24454 ops/sec)"
        heartbeat_pattern = re.compile(r"tried (\d+) keys \((\d+) ops/sec\)")
        # Regex for parsing found keys
        found_pattern = re.compile(r"FOUND: (.*)")

        with open(self.log_file, "a") as log:
            for line in iter(self.process.stdout.readline, ""):
                line = line.strip()
                if not line:
                    continue

                log.write(f"{time.strftime('%Y-%m-%d %H:%M:%S')} - {line}\n")
                log.flush()

                with self.lock:
                    self.status["last_heartbeat"] = line

                    hb_match = heartbeat_pattern.search(line)
                    if hb_match:
                        self.status["total_tried"] = int(hb_match.group(1))
                        self.status["iops"] = int(hb_match.group(2))

                    found_match = found_pattern.search(line)
                    if found_match:
                        self.status["found_keys"].append({
                            "timestamp": time.strftime('%Y-%m-%d %H:%M:%S'),
                            "key_info": found_match.group(1)
                        })
                        self._save_found_key(found_match.group(1))

            self.status["running"] = False

    def _save_found_key(self, key_info):
        with open("found_keys.json", "a") as f:
            f.write(json.dumps({
                "timestamp": time.strftime('%Y-%m-%d %H:%M:%S'),
                "key_info": key_info
            }) + "\n")

    def get_status(self):
        with self.lock:
            return self.status.copy()
