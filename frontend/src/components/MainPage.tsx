import {useState} from "react";
import BasicAppBar from "./BasicAppBar";
import GroupForm from "./GroupForm";
import CenteredTabs from "./CenteredTabs";
import EmailVerificationForm from "./EmailVerficationForm";

const MainPage = () => {
    const [tabValue, setTabValue] = useState<number>(0);
    const [emailVerified, setEmailVerified] = useState<boolean>(false);

    const handleChange = (event: React.SyntheticEvent, newValue: number) => {
        if (newValue === 0 || emailVerified) {
            setTabValue(newValue);
        }
    };

    const handleVerifyEmail = () => {
        console.log("set to verified")
        setEmailVerified(true);
        setTabValue(1);
    };

    return (
        <div style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
            <div style={{ flex: "0 0 auto" }}>
                <BasicAppBar />
            </div>
            <div style={{ flex: "1 1 auto", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", padding: "2rem" }}>
                <CenteredTabs value={tabValue} onChange={handleChange} emailVerified={emailVerified} />
                <div style={{ marginTop: "1rem", width: "100%" }}>
                    {tabValue === 0 && !emailVerified && (
                        <EmailVerificationForm
                            onVerify={handleVerifyEmail}
                        />
                    )}
                    {(tabValue === 1 || tabValue === 2 || emailVerified) && (
                        <GroupForm 
                            groupIdentifier={tabValue === 1 ? "groupName" : "groupId"}
                            label={tabValue === 1 ? "Group Name" : "Group ID"} 
                        />
                    )}
                </div>
            </div>
        </div>
    );
};
export default MainPage