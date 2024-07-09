import React from "react";
import { Button, TextField, Grid, Paper, Typography } from "@mui/material";
import { useMutation } from "@tanstack/react-query";

import submitVerifyEmail from "../queries/submitVerifyEmail";

interface EmailVerificationFormProps {
    onVerify: () => void;
}

const EmailVerificationForm = ({ onVerify }: EmailVerificationFormProps) => {
    const submitVerifyEmailMutation = useMutation(submitVerifyEmail, {
        onSuccess: (data) => {
            onVerify();
            // todo 
            console.log(data)
        },
        onError: (error: Error) => {
            console.log(error.message);
        },
    });

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget)
        const email = formData.get("email") as string;
        submitVerifyEmailMutation.mutate({email})
    };

    return (
        <Grid container justifyContent="center" alignItems="center">
            <Grid item xs={10} sm={8} md={6} lg={4}>
                <Paper elevation={3} style={{ padding: "2rem" }}>
                    <Typography variant="h5" gutterBottom>
                        Verify Your Email
                    </Typography>
                    <form onSubmit={handleSubmit}>
                        <Grid container spacing={2}>
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
                            <Grid item xs={6}>
                                <Button type="submit" variant="contained" color="primary">
                                    Verify
                                </Button>
                            </Grid>
                            <Grid item xs={6}>
                                <Button variant="outlined" color="primary" onClick={onVerify}>
                                    Already Verified
                                </Button>
                            </Grid>
                        </Grid>
                    </form>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default EmailVerificationForm;