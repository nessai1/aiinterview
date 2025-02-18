import React from "react";
import {CountdownTimer} from "@/components/internal/chat/CountdownTimer.tsx";


type TProps = {
    interviewId: string
}

const Chat: React.FC<TProps> = (props: TProps) => {

    return (
        <>
            <CountdownTimer seconds={180}/>
            Interview: {props.interviewId}  Interview: {props.interviewId}  Interview: {props.interviewId}
        </>
    );
}

export default Chat;