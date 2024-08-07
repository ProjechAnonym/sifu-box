import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { Provider } from "react-redux";
import { store } from "@/redux/store";
import { Toaster, ToastBar, toast } from "react-hot-toast";
import UIProvider from "@/uiProvider";
import App from "@/App";
import "@/index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <UIProvider>
          <Toaster>
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
        </UIProvider>
      </BrowserRouter>
    </Provider>
  </React.StrictMode>
);
