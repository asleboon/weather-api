import * as resources from '@pulumi/azure-native/resources';
import * as web from "@pulumi/azure-native/web";
import * as pulumi from "@pulumi/pulumi";

// Create an Azure Resource Group
const resourceGroup = new resources.ResourceGroup('weather-api-dev', {
	location: 'West Europe',
});

const plan = new web.AppServicePlan('free-plan', {
	resourceGroupName: resourceGroup.name,
	location: resourceGroup.location,
	kind: 'Linux',
	reserved: true, // Creates Windows machine if not set.
	sku: {
		name: 'F1',
		tier: 'Free'
	}
});

const imageInDockerHub = 'asleboon/weather-api:latest';

const app = new web.WebApp("weather-api-dev", {
    resourceGroupName: resourceGroup.name,
	location: resourceGroup.location,
    serverFarmId: plan.id,
    siteConfig: {
		appSettings: [
            {
                name: "WEBSITES_ENABLE_APP_SERVICE_STORAGE",
                value: "false",
            },
            {
                name: "WEBSITES_PORT",
                value: "8080",
            },
        ],
        alwaysOn: false,
        linuxFxVersion: `DOCKER|${imageInDockerHub}`,
    },
    httpsOnly: true,
});

export const appUrl = pulumi.interpolate`https://${app.defaultHostName}`;

// docker login --username=asleboon --email=asle.berge@gmail.com
// docker build -t weather-api-app
// docker tag weather-api-app:latest asleboon/weather-api:latest
// docker push asleboon/weather-api:latest

// pulumi stack output appUrl
// weather-api-dev49f89349.azurewebsites.net
//

// Ensure compaitbility or it will not start
//docker buildx create --use
// docker buildx build --platform linux/amd64 -t yourdockerhubusername/yourdockerhubrepository:yourtag .


