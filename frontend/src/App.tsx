import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import MainPage from "./components/MainPage";
import GroupDetailsPage from "./components/GroupDetailsPage";

const queryClient = new QueryClient();

const App = () => {
    return (
        <QueryClientProvider client={queryClient}>
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<MainPage />} />
                    <Route path="/group-details/:id" element={<GroupDetailsPage />} />
                </Routes>
            </BrowserRouter>
        </QueryClientProvider>
    );
};

const container = document.getElementById("root")!;
const root = createRoot(container);
root.render(<App />);