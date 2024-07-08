import React from 'react';
import BasicAppBar from "./BasicAppBar";
import CreateGroupForm from "./CreateGroupForm";
import JoinGroupForm from './JoinGroupForm';
import CenteredTabs from './CenteredTabs';

const MainPage = () => {
    const [tabValue, setTabValue] = React.useState(0);

    const handleChange = (event, newValue) => {
        setTabValue(newValue);
    };

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            <div style={{ flex: '0 0 auto' }}>
                <BasicAppBar />
            </div>
            <div style={{ flex: '0 0 auto' }}>
                <CenteredTabs value={tabValue} onChange={handleChange}/>
            </div>
            <div style={{ flex: '1 1 auto', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                {tabValue === 0 ? <CreateGroupForm /> : <JoinGroupForm />}
            </div>
        </div>
    );
};
export default MainPage