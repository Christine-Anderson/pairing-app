import BasicAppBar from './BasicAppBar';
import Typography from '@mui/material/Typography';
import NameList from './NameList';

const GroupDetailsPage = () => {
    const groupName = "My Group"
    const names = ["Jane", "John", "Alice", "Bob"];

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            <div style={{ flex: '0 0 auto' }}>
                <BasicAppBar />
            </div>
            <div style={{ flex: '1 1 auto', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', padding: '2rem' }}>
                <Typography variant="h4" align="center" gutterBottom>
                    {groupName}
                </Typography>
                <div style={{ marginTop: '1rem', width: '100%' }}>
                    <NameList names={names}/>
                </div>
            </div>
        </div>
    );
};
export default GroupDetailsPage