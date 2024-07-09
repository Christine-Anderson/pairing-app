import { TextField, Button, Grid, Paper, Typography } from "@mui/material";

interface GroupFormProps {
    groupIdentifier: string;
    label: string;
}
const GroupForm = ({groupIdentifier, label}: GroupFormProps) => {
    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget)
        const name = formData.get("name");
        const email = formData.get("email");
        const group = formData.get(groupIdentifier);

        console.log('Form Data:', { name, email, group });
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