import React from 'react';
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import BasicAppBar from "./components/BasicAppBar";
import CreateGroupForm from "./components/CreateGroupForm";
import JoinGroupForm from './components/JoinGroupForm';
import CenteredTabs from './components/CenteredTabs';

const queryClient = new QueryClient();

const App = () => {
    const [tabValue, setTabValue] = React.useState(0);

    const handleChange = (event, newValue) => {
        setTabValue(newValue);
    };

    return (
        <QueryClientProvider client={queryClient}>
            <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
                <div style={{ flex: '0 0 auto' }}>
                    <BasicAppBar />
                </div>
                <div style={{ flex: '0 0 auto' }}>
                    <CenteredTabs value={tabValue} onChange={handleChange}/>
                </div>
                <div style={{ flex: '1 1 auto', display: 'flex', justifyContent: 'center', alignItems: 'center', padding: '16px' }}>
                    {tabValue === 0 ? <CreateGroupForm /> : <JoinGroupForm />}
                </div>
            </div>
        </QueryClientProvider>
    );
};

const container = document.getElementById("root");
const root = createRoot(container);
root.render(<App />);