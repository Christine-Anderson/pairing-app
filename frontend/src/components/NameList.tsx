import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import Button from "@mui/material/Button";
import Paper from "@mui/material/Paper";
import FormControl from "@mui/material/FormControl";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import Typography from "@mui/material/Typography";
import LinearProgress from '@mui/material/LinearProgress';
import Alert from '@mui/material/Alert';
import AlertTitle from '@mui/material/AlertTitle';
import { GroupMember } from "../types";
import submitGenerateAssignments from "../queries/submitGenerateAssignments";

interface NameListProps {
    groupId: string;
    jwt: string;
    groupMembers: GroupMember[];
    isLoading: boolean;
    onGenerateAssignments: () => void;
}

const NameList = ({groupId, jwt, groupMembers, isLoading, onGenerateAssignments}: NameListProps) => {
    const [selectedValues, setSelectedValues] = useState<{ [key: string]: string[] }>({});

    const submitGenerateAssignmentsMutation = useMutation(submitGenerateAssignments, {
        onSuccess: (data) => {
            onGenerateAssignments();
            console.log(data);
        },
        onError: (error: Error) => {
            <Alert severity="error">
                <AlertTitle>Error</AlertTitle>
                {error.message}
            </Alert>
        },
    });

    const handleChange = (memberId: string, event: React.ChangeEvent<{ value: unknown }>) => {
        const selectedIds = event.target.value as string[];
        setSelectedValues(prev => ({
            ...prev,
            [memberId]: selectedIds,
        }));
    };

    const handleButtonClick = () => {
        submitGenerateAssignmentsMutation.mutate({groupId, jwt, restrictions: selectedValues});
        console.log("Selected IDs:", selectedValues);
    };

    return (
        <Paper elevation={3} style={{ padding: "1rem", maxWidth: "20rem", margin: "auto" }}>
            { isLoading && <LinearProgress/> }
            {
                groupMembers && 
                <Typography variant="subtitle1" align="center" gutterBottom>
                    Please indicate any restrictions in who the participants can be paired with.
                </Typography>
            }
            <List>
                {groupMembers && groupMembers.map(({memberId, name}, index) => (
                    <ListItem key={index}>
                        <ListItemText primary={name} />
                        <FormControl style={{ minWidth: "12rem", marginLeft: "1rem" }}>
                            <Select
                                multiple
                                value={selectedValues[memberId] || []}
                                onChange={(event) => handleChange(memberId, event as React.ChangeEvent<HTMLInputElement>)}
                                renderValue={(selected) => {
                                    const selectedNames = (selected).map(memberId => {
                                        const selectedPerson = groupMembers.find(item => item.memberId === memberId);
                                        return selectedPerson ? selectedPerson.name : "";
                                    });
                                    return selectedNames.join(", ");
                                }}
                                style={{ minWidth: "9em" }}
                            >
                                {groupMembers
                                    .filter(member => member.memberId !== memberId)
                                    .map(({ memberId, name }) => (
                                        <MenuItem key={memberId} value={memberId}>
                                            {name}
                                        </MenuItem>
                                    ))}
                            </Select>
                        </FormControl>
                    </ListItem>
                ))}
            </List>
            <div style={{ textAlign: "center", marginTop: "1em" }}>
                <Button variant="contained" color="primary" onClick={handleButtonClick}>
                    Generate Assignments
                </Button>
            </div>
        </Paper>
    );
};
export default NameList;