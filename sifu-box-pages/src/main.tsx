import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { Toaster, ToastBar, toast } from "react-hot-toast";
import { Provider } from "react-redux";
import { store } from "./redux/store";
import App from "./App";
import { NavProvider } from "./provider";
import { toast_config } from "./config/toast";
import "bootstrap-icons/font/bootstrap-icons.css";
import "@/styles/globals.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <NavProvider>
          <Toaster
            gutter={8}
            toastOptions={toast_config}
          >
            {(t) => (
              <ToastBar toast={t}>
                {({ icon, message }) => (
                  <>
                    {icon}
                    {message}
                    {t.type !== "loading" && (
                      <button onClick={() => toast.dismiss(t.id)}>
                        <i className="bi bi-x-lg" />
                      </button>
                    )}
                  </>
                )}
              </ToastBar>
            )}
          </Toaster>
          <App />
        </NavProvider>
      </BrowserRouter>
    </Provider>
  </React.StrictMode>,
);
