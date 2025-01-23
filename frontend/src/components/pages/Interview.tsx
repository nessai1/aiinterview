import React from "react";
import {useParams} from "react-router-dom";



const Interview: React.FC = () => {

    const { interviewId } = useParams();

    return (
        <>Interview: {interviewId} </>
    );
}

export default Interview;