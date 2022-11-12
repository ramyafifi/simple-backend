# Octank University CICD Pipeline

This contains the artifacts from my AWSome Builder Project (May 2021).  The AB3 project is made up of two individual code projects, a Serverless VueJS Frontend and a simple Go app that is utilized as an API.

Contact/Author:<br>
David Stielstra, *AWS WWPS Sr. Solutions Architect* <br>
dstiel@amazon.com

### Version 1.0.0

## Overview
Octank University is a R1 university that has an internal development team.  The development team has been looking into utilizing containers for their applications but first wanted to modernize their development practices.

## Challenge

While currently using a git based code repository when it comes to building and deployment of changes the team has been doing this manually.  They know there is a better way and would like to see how the process may be improved through the utilization of the AWS Code * services.

## Solution
A combination of CodeCommit, CodeBuild, Code Deploy orchestrated by Code Pipeline are utilized to handle releases of the simple Vue.js frontend and GO application running on AWS ECS Fargate.

This project is part B - A Simple GO Docker backend API project that is running in ECS Fargate.  Deployed using CICD pipeline utilizing CodePipeline built from CodeCommit, CodeBuild and CodeDeploy.  The GO app simply retrieves the task metadata to be returned to the caller.

The buildspec file in this solution is constructed in 3 phases to log in to the Amazon ECR, build/tag the Docker image, and deploy the Docker image to ECR with 3 artifact files used by Fargate (imageDetail.json, taskdef.json and appspec.yml).

## Architecture

![See architecture/Architecture.png for diagram](architecture/CICD Application.png)
![See architecture/CICD Pipeline.png for diagram](architecture/CICD Pipeline.png)

## Deployment Notes
This was created using the AWS console. It would be nice to create CloudFormation templates as the next iteration.  You should have already followed the steps outlined in the Serverless VueJS Frontend project.  There should be a CodeCommit repository populated with the code from this Simple Go app solution.

To Deploy:
1.) Create an ECR Repository to store the Docker images that will be built.

2.) Create an Amazon ECS Cluster as a Networking Only clusted as you'll be using AWS Fargate, creating a Service and Tasks in the next steps.

3.) Create a new Task Definition with the launch type of 'Fargate'
- You will need to create a IAM Role for ecsTaskExecution with the 'AmazonECSTaskExecutionRolePolicy' attached.
- Network mode should be awsvpc
- Operating System = Linux
- 0.5GB and 0.25 vCPU should be sufficient
- Add a container - with the image that is pulled from your ECR repository - you may need to manually push an image to ECR.


4.) In the newly created ECS Cluster create a new Service with a launch type of Fargate.
- Select the Task Definition created in the previous step.
- You may find it useful to start with a quantity of 2 tasks so you can see on the front end the AZ value update.
- For deployments choose Blue/Green
- When it comes to Load Balancing choose the Application Load Balancer selecting the ALB created in Part A, be sure to Add the container to the load balancer, utilizing the Target Group.

5.) Create CodePipeline to handle API changes
- Source Stage: Source should be CodeCommit using the repo/branch for the API.
- Use Amazon CloudWatch Events and the Full clone output artifact format.
- SourceVariables and SourceArtifact for variable namespace and output artifacts respectively.

- Build Stage:  Build Action - use AWS CodeBuild
- Input artifact should be SourceArtifact
- Create a build project that uses the standard codebuild image and make sure to use a buildspecfile giving the filename 'buildspec.yml'.  Be sure to add environment variables of 'AWS_DEFAULT_REGION', 'AWS_ACCOUNT_ID', 'IMAGE_TAG', 'IMAGE_REPO_NAME' with values such as 'us-east-1', '<your aws account number>','latest','ab3' .
- Build type = 'Single Build'
- Variable namespace - 'BuildVariables'
- Output artifacts - 'BuildArtifact'

- Deploy Stage: Manual approval to ECS Blue/Green deployment.
- Manual approval should utilize SES so you will need to setup a SNS Topic with your email address to receive the notification.
- Blue/Green deployment will utilize the Action Provider of Amazon ECS (Blue/Green)
- Input artifact - 'BuildArtifact'
- Create both a AWS CodeDeploy application and CodeDeploy deployment group
- ECS Task definition - BuildArtifact = 'taskdef.json'
- ApsSpec file - BuildArtifact = 'appspec.yml'
- Input artifact with image details = BuildArtifact
- Placeholder text in the task definition - 'IMAGE1_NAME'
- Variable namespace - 'DeployVariables'

6.) Make a code change and commit the change to your api repo.  This should trigger the CodePipeline execution.  You should receive an email notification that a new deployment is ready for review.  You'll need to log into the console to approve the Deploy stage and also the blue/green deployment.