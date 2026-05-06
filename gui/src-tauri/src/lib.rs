use tauri_plugin_shell::ShellExt;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_shell::init())
        .setup(|app| {
            let shell = app.shell();
            let sidecar = shell.sidecar("pwny-server").unwrap();
            
            tauri::async_runtime::spawn(async move {
                let (mut _rx, mut _child) = sidecar.spawn().expect("Failed to spawn sidecar");
                // We could monitor output here if needed
            });
            
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
