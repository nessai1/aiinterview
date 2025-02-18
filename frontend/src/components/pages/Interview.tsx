import React from "react";
import {useParams} from "react-router-dom";
import Chat from "@/components/internal/chat/Chat.tsx";



const Interview: React.FC = () => {

    const { interviewId } = useParams();
    if (interviewId === undefined) {
        window.location.href = '/';
        return;
    }

    return (
        <><Chat interviewId={interviewId} /></>
    );
}

export default Interview;