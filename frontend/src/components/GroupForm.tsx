import { TextField, Button, Grid, Paper, Typography } from "@mui/material";
import { useMutation } from "@tanstack/react-query";
import Alert from '@mui/material/Alert';
import AlertTitle from '@mui/material/AlertTitle';

import submitCreateGroup from "../queries/submitCreateGroup";
import submitJoinGroup from "../queries/submitJoinGroup";

interface GroupFormProps {
    groupIdentifier: string;
    label: string;
    onSuccess: () => void;
}

const GroupForm = ({groupIdentifier, label, onSuccess}: GroupFormProps) => {
    const submitCreateGroupMutation = useMutation(submitCreateGroup, {
        onSuccess: (data) => {
            onSuccess();
            console.log(data);
        },
        onError: (error: Error) => {
            <Alert severity="error">
                <AlertTitle>Error</AlertTitle>
                {error.message}
            </Alert>
        },
    });

    const submitJoinGroupMutation = useMutation(submitJoinGroup, {
        onSuccess: (data) => {
            onSuccess();
            console.log(data);
        },
        onError: (error: Error) => {
            <Alert severity="error">
                <AlertTitle>Error</AlertTitle>
                {error.message}
            </Alert>
        },
    });

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        const name = formData.get("name") as string;
        const email = formData.get("email") as string;
        const group = formData.get(groupIdentifier) as string;
        
        if (groupIdentifier === "groupName") {
            submitCreateGroupMutation.mutate({ name, email, groupName: group });
        } else if (groupIdentifier === "groupId") {
            submitJoinGroupMutation.mutate({ name, email, groupId: group })
        }

        event.currentTarget.reset();
    };

    return (
        <Grid container justifyContent="center" alignItems="center">
            <Grid item xs={10} sm={8} md={6} lg={4}>
                <Paper elevation={3} style={{ padding: "2rem" }}>
                    <Typography variant="h5" gutterBottom>
                        Create a Group
                    </Typography>
                    <form onSubmit={handleSubmit}>
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    required
                                    id="name"
                                    name="name"
                                    label="Name"
                                    variant="outlined"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    required
                                    id="email"
                                    name="email"
                                    label="Email"
                                    variant="outlined"
                                    type="email"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    required
                                    id={groupIdentifier}
                                    name={groupIdentifier}
                                    label={label}
                                    variant="outlined"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <Button type="submit" variant="contained" color="primary">
                                    Submit
                                </Button>
                            </Grid>
                        </Grid>
                    </form>
                </Paper>
            </Grid>
        </Grid>
    );
}

export default GroupForm;