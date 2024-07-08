import React from 'react';
import BasicAppBar from "./BasicAppBar";
import CreateGroupForm from "./CreateGroupForm";
import JoinGroupForm from './JoinGroupForm';
import CenteredTabs from './CenteredTabs';

const MainPage = () => {
    const [tabValue, setTabValue] = React.useState<number>(0);

    const handleChange = (event: React.SyntheticEvent, newValue: number) => {
        setTabValue(newValue);
    };

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            <div style={{ flex: '0 0 auto' }}>
                <BasicAppBar />
            </div>
            <div style={{ flex: '1 1 auto', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', padding: '2rem' }}>
                <CenteredTabs value={tabValue} onChange={handleChange}/>
                <div style={{ marginTop: '1rem', width: '100%' }}>
                    {tabValue === 0 ? <CreateGroupForm /> : <JoinGroupForm />}
                </div>
            </div>
        </div>
    );
};
export default MainPage