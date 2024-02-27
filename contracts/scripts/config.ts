export interface Config {
    runs:         number;
    dst:          string;
    submit:       boolean;
    Logger:       Logger;
    Ethereum:     Ethereum;
    BurnCircuit:  Circuit;
    ClaimCircuit: Circuit;
}

export interface Circuit {
    constraintSystemPath: string;
    provingKeyPath:       string;
    verifyingKeyPath:     string;
}

export interface Ethereum {
    host:                  string;
    privateKey:            string;
    rollupContract:        string;
    oracleMockContract:    string;
    burnVerifierContract:  string;
    claimVerifierContract: string;
}

export interface Logger {
    LogLevel: string;
    Format:   string;
}
