import React from 'react';
import BasicAppBar from './BasicAppBar';
import Typography from '@mui/material/Typography';
import NameList from './NameList';

const GroupDetailsPage = () => {
    const groupName = "Example Group"
    const names = ["Jane", "John", "Alice", "Bob"];

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            <div style={{ flex: '0 0 auto' }}>
                <BasicAppBar />
            </div>
            <div style={{ flex: '0 0 auto', padding: '2rem' }}>
                <Typography variant="h4" align="center" gutterBottom>
                    {groupName}
                </Typography>
                <Typography variant="subtitle1" align="center" gutterBottom>
                    Put an explanation here.
                </Typography>
            </div>
            <div style={{ flex: '1 1 auto', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <NameList names={names}/>
            </div>
        </div>
    );
};
export default GroupDetailsPage