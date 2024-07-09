import { Navigate, useParams, useLocation } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import queryString from "query-string";
import BasicAppBar from "./BasicAppBar";
import Typography from "@mui/material/Typography";
import NameList from "./NameList";
import fetchGroupDetails from "../queries/fetchGroupDetails";

const GroupDetailsPage = () => {
    const {groupId} = useParams<{ groupId: string }>();
    const location = useLocation();
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
    const { jwt } = queryString.parse(location.search) as { jwt: string };
    const { data: groupDetails, isLoading, isError } = useQuery(
        ["groupDetails", groupId],
        () => fetchGroupDetails({ groupId: groupId || "", jwt })
    );

    if (!jwt || isError) {
        return <Navigate to="/" />;
    }

    const groupName = groupDetails?.groupName || "My Group";
    const groupMembers = groupDetails?.groupMembers || [];

    return (
        <div style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
            <div style={{ flex: "0 0 auto" }}>
                <BasicAppBar />
            </div>
            <div style={{ flex: "1 1 auto", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", padding: "2rem" }}>
                <Typography variant="h4" align="center" gutterBottom>
                    {groupName}
                </Typography>
                <div style={{ marginTop: "1rem", width: "100%" }}>
                    <NameList groupMembers={groupMembers} isLoading={isLoading}/>
                </div>
            </div>
        </div>
    );
};
export default GroupDetailsPage