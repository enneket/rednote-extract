#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::net::TcpStream;
use std::path::PathBuf;
use std::process::{Child, Command, Stdio};
use std::sync::Mutex;
use std::time::Duration;

struct AppState {
    gateway: Mutex<Option<Child>>,
}

impl Drop for AppState {
    fn drop(&mut self) {
        if let Some(mut child) = self.gateway.lock().unwrap().take() {
            let _ = child.kill();
        }
    }
}

fn get_gateway_path() -> Option<PathBuf> {
    let exe_dir = std::env::current_exe().ok()?.parent()?.to_path_buf();

    for name in ["xhs-gateway.exe", "xhs-gateway"] {
        let path = exe_dir.join(name);
        if path.exists() {
            return Some(path);
        }
    }

    exe_dir
        .parent()
        .and_then(|p| {
            let fp = p.join("gateway").join("xhs-gateway.exe");
            if fp.exists() {
                Some(fp)
            } else {
                None
            }
        })
}

fn wait_gateway_ready() {
    for _ in 0..20 {
        if TcpStream::connect("127.0.0.1:8000").is_ok() {
            println!("Gateway ready at http://127.0.0.1:8000");
            return;
        }
        std::thread::sleep(Duration::from_millis(500));
    }
    println!("Gateway started (health check skipped)");
}

#[cfg(windows)]
fn spawn_gateway(path: &PathBuf) -> Option<Child> {
    use std::os::windows::process::CommandExt;
    const CREATE_NO_WINDOW: u32 = 0x08000000;
    Command::new(path)
        .creation_flags(CREATE_NO_WINDOW)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .ok()
}

#[cfg(not(windows))]
fn spawn_gateway(path: &PathBuf) -> Option<Child> {
    Command::new(path)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .ok()
}

fn main() {
    let gateway = get_gateway_path().and_then(|p| {
        println!("Found gateway at {:?}", p);
        spawn_gateway(&p)
    });

    if gateway.is_some() {
        wait_gateway_ready();
    } else {
        println!("Gateway exe not found, skipping auto-start");
    }

    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .manage(AppState {
            gateway: Mutex::new(gateway),
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
