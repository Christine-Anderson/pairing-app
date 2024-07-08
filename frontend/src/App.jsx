import React from 'react';
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

const App = () => {
    return (
        <QueryClientProvider client={queryClient}>
            <h1>Hello World</h1>
        </QueryClientProvider>
    );
};

const container = document.getElementById("root");
const root = createRoot(container);
root.render(<App />);