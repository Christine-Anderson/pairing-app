import Box from "@mui/material/Box";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";

interface CenteredTabsProps {
    value: number;
    onChange: (event: React.SyntheticEvent, newValue: number) => void;
    emailVerified: boolean;
}

const CenteredTabs = ({value, onChange, emailVerified}: CenteredTabsProps) => {
    return (
        <Box sx={{ width: "100%", bgcolor: "background.paper" }}>
            <Tabs value={value} onChange={onChange} centered>
                <Tab label="Verify Email" disabled={emailVerified} />
                <Tab label="Create Group" disabled={!emailVerified} />
                <Tab label="Join Group" disabled={!emailVerified} />
            </Tabs>
        </Box>
    );
}

export default CenteredTabs;