import Box from '@mui/material/Box';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';

interface CenteredTabsProps {
    value: number;
    onChange: (event: React.SyntheticEvent, newValue: number) => void;
}

const CenteredTabs = ({value, onChange}: CenteredTabsProps) => {
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