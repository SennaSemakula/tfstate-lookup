package cmd

import (
	"github.com/SennaSemakula/tfstate-lookup/pkg/terraform"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	account    string
	bucketName string
	rootCmd    = &cobra.Command{
		Use:   "tfstate-lookup",
		Short: "Check what terraform infra deployed on AWS",
		Long: `CLI tool to query what terraform resources are deployed in an AWS account. Can be used to query any AWS account.
		
		Example:
		tfstate-lookup --account <account> --bucket_name <bucket>
		`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Checking to see what resources are deployed in %s AWS account...\n", account)
			err := terraform.GetAWSResources(bucketName)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&account, "account", "", "AWS account to query terraform resources")
	rootCmd.PersistentFlags().StringVar(&bucketName, "bucket_name", "", "AWS bucket name where terraform state files are stored")
	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if account == "" || bucketName == "" {
		log.Fatal("flags account and bucket account must be set")
	}
	os.Setenv("AWS_PROFILE", account)
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
}
