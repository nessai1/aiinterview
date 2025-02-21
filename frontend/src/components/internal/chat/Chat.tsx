import React from "react";
import {CountdownTimer} from "@/components/internal/chat/CountdownTimer.tsx";
import SectionBlock from "@/components/internal/chat/SectionBlock.tsx";


type TProps = {
    interviewId: string
}

const Chat: React.FC<TProps> = (props: TProps) => {

    return (
        <>
            {/*<CountdownTimer seconds={180}/>*/}
            {/*Interview: {props.interviewId}  Interview: {props.interviewId}  Interview: {props.interviewId}*/}
            <SectionBlock section={{}} />
        </>
    );
}

export default Chat;