import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { Toaster, ToastBar, toast } from "react-hot-toast";
import { Provider } from "react-redux";
import { store } from "@/redux/store";
import { NavProvider } from "./provider.tsx";
import App from "./App.tsx";
import "bootstrap-icons/font/bootstrap-icons.css";
import "@/styles/globals.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <NavProvider>
          <Toaster
            gutter={8}
            toastOptions={{
              // Define default options
              className: "",
              duration: 5000,
              style: {
                background: "#F0E8EE",
                color: "#4C4C4C",
                fontWeight: "bolder",
                fontSize: "0.8rem",
              },
              success: {
                style: {
                  background: "#86efac",
                  color: "#5c5c5c",
                  fontWeight: "bolder",
                },
              },
              error: {
                style: {
                  background: "#ff607e",
                  color: "#ededed",
                  fontWeight: "bolder",
                },
              },
            }}
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
  </React.StrictMode>
);
