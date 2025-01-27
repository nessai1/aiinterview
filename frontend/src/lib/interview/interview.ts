type Interview = {
    title: string,
    complete: boolean,
    timing: number,
    topics: Topic[]
}

type Topic = {
    name: string
    grade: Grade
}

enum Grade {
    Junior = 'junior',
    Middle = 'middle',
    Senior = 'senior'
}

export {
    Interview,
    Topic,
    Grade
}
