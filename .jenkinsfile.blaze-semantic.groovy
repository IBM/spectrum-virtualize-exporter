def buildInfo = bzSemantic(
    pushLatestTag              : true, 
    architecture               : 's390x',
    buildDockerImage           : false,
    useAnalysisCodeCoverage    : false,
    useAnalysisDependencyCheck : false,
    useAnalysisOsscCheck       : false,
    addDependenciesToBuildInfo : false,
    buildNodeLabel: 'blaze-build-agent-s390-colo-new',
    baseVersion : '0.11.4',
    verbose: true,
    useAnalysisSonarQube: true,
    sonarQube: [
	    useEE: true
	  ],
    native: [
        kind: "go",
        build: [
            platforms: ["linux/arm64", "linux/s390x"]
        ]
    ]
)