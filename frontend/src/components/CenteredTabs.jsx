import * as React from 'react';
import Box from '@mui/material/Box';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';

const CenteredTabs = ({value, onChange}) => {
    return (
        <Box sx={{ width: '100%', bgcolor: 'background.paper' }}>
            <Tabs value={value} onChange={onChange} centered>
                <Tab label="Create Group" />
                <Tab label="Join Group" />
            </Tabs>
        </Box>
    );
}

export default CenteredTabs;