//
//  NotificationPublishViewController.m
//  Squirrel
//
//  Created by JamesMao on 4/12/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "NotificationPublishViewController.h"
#import "AFNetworking.h"
#import "Define.h"
#import "MeViewData.h"

@interface NotificationPublishViewController ()
{
    IBOutlet UITextField* classNameTextField_;
    IBOutlet UITextField* classTimeTextField_;
    IBOutlet UITextField* classTeachTextField_;
    IBOutlet UITextField* classDescriptionTextField_;
    IBOutlet UITextField* classMaxStudentNumberTextField_;
    
    IBOutlet UIButton* deleteButton_;
    IBOutlet UIButton* registerButton_;
    
    NotificationDataItem* notificationDataItem_;
    UserRegisterDataItem* userRegisterDataItem_;
    
    EditType editType_;
}
- (IBAction)onDelete:(id)sender;
- (IBAction)onRegister:(id)sender;

@end


@implementation NotificationPublishViewController

- (id)init
{
    if (self = [super init])
    {
        notificationDataItem_ = [[NotificationDataItem alloc] init];
        
        userRegisterDataItem_ = [[UserRegisterDataItem alloc] init];
        
        self.hidesBottomBarWhenPushed = YES;
        editType_ = EditType_Add;
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    self.navigationItem.rightBarButtonItem = [[UIBarButtonItem alloc] initWithTitle:@"Save"
                                                                              style:UIBarButtonItemStyleDone
                                                                             target:self
                                                                             action:@selector(onSave)];
}

- (void)onSave
{
    [notificationDataItem_ setClassName: classNameTextField_.text];
    [notificationDataItem_ setClassTime: classTimeTextField_.text];
    [notificationDataItem_ setClassTeacher: classTeachTextField_.text];
    [notificationDataItem_ setClassDescription: classDescriptionTextField_.text];
    [notificationDataItem_ setClassMaxStudent: classMaxStudentNumberTextField_.text];
    
    [[AdminPublishViewData singleton] setNotificationDataItem:notificationDataItem_];
    
    [self postWithAFHttp];
}

- (void)postWithAFHttp
{
    AFHTTPRequestOperationManager *manager = [AFHTTPRequestOperationManager manager];

    NSDictionary *parameters = [notificationDataItem_ getDataItemEelements];
    
    NSString* urlString = [NSString stringWithFormat:@"%@/API1/Post", SERVER_IP];
    [manager POST:urlString parameters:parameters success:^(AFHTTPRequestOperation *operation, id responseObject) {
        NSLog(@"JSON: %@", responseObject);
        [self.navigationController popViewControllerAnimated:YES];
    } failure:^(AFHTTPRequestOperation *operation, NSError *error) {
        NSLog(@"Error: %@", error);
    }];
}

- (void)viewWillAppear:(BOOL)animated
{
    if (EditType_View == editType_
        || EditType_Editable == editType_)
    {
        dispatch_async(BackGroundQueue, ^{
            
            NSString* key = [notificationDataItem_ getUniqueKey];
            NSLog(@"%@", key);
            NSString* urlString = [NSString stringWithFormat:@"%@%@", OneClassURLBase, key];
            NSLog(@"%@", urlString);
            NSData* data = [NSData dataWithContentsOfURL:[NSURL URLWithString: urlString]];
            if (data == nil)
            {
                NSLog(@"data is nil");
                return;
            }
            [self performSelectorOnMainThread:@selector(fetchedOneClass:) withObject:data waitUntilDone:YES];
        });
    }
    
    [self updateEditType];
    [self updateUIFromData];
}

- (void)fetchedOneClass:(NSData *)responseData
{
    if (responseData == nil)
    {
        return;
    }
    
    //parse out the json data
    NSError* error;
    NSDictionary* json = [NSJSONSerialization JSONObjectWithData:responseData //1
                                                         options:kNilOptions
                                                           error:&error];
    NSLog(@"dictionary data %@",json);
    
    NSMutableDictionary* mutableDictionary = [[NSMutableDictionary alloc] init];
    [mutableDictionary setDictionary:json];
    [notificationDataItem_ setDataItemEelements:mutableDictionary];
    
    [self updateEditType];
    [self updateUIFromData];
}


- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (void)setEditType:(EditType)type
{
    editType_ = type;
}

- (void)updateEditType
{
    BOOL editable = YES;
    BOOL deletable = YES;
    BOOL registable = NO;
    switch (editType_)
    {
        case EditType_View:
        {
            editable = NO;
            deletable = NO;
            registable = YES;
        }
        break;
        case EditType_Editable:
        {
            editable = YES;
        }
        break;
        case EditType_Add:
        {
            editable = YES;
            deletable = NO;
        }
        break;
            
        default:
            break;
    }
    [self setEditable:editable];
    
    deleteButton_.hidden = deletable == NO;
    registerButton_.hidden = registable == NO;
    self.navigationItem.rightBarButtonItem.enabled = editable;
}

- (void)setNotificationDataItem:(NotificationDataItem*)item
{
    notificationDataItem_ = item;
}

- (IBAction)onDelete:(id)sender
{
    [self postWithDeleteMessage];
    [self.navigationController popViewControllerAnimated:YES];
}


- (void)postWithDeleteMessage
{
    AFHTTPRequestOperationManager *manager = [AFHTTPRequestOperationManager manager];
    
    NSString* uniqueKey = [notificationDataItem_ getUniqueKey];
    
    NSString* urlString = [NSString stringWithFormat:@"%@%@%@", SERVER_IP, URL_PART_DELETE_ONE_CLASS, uniqueKey];
    [manager POST:urlString parameters:nil success:^(AFHTTPRequestOperation *operation, id responseObject) {
        NSLog(@"JSON: %@", responseObject);
    } failure:^(AFHTTPRequestOperation *operation, NSError *error) {
        NSLog(@"Error: %@", error);
    }];
}

- (IBAction)onRegister:(id)sender
{
    [self postWithRegisterMessage];
    [self.navigationController popViewControllerAnimated:YES];
}

- (void)postWithRegisterMessage
{
    
    NSString* uniqueNameKey = [MeViewData singleton].userUniqueName;
    [userRegisterDataItem_ setUniqueKey:[notificationDataItem_ getUniqueKey]];
    [userRegisterDataItem_ setClassName:[notificationDataItem_ getClassName]];
    
    [userRegisterDataItem_ setUserName:uniqueNameKey];
    [userRegisterDataItem_ setUserUniqueKey:uniqueNameKey];
    
    AFHTTPRequestOperationManager *manager = [AFHTTPRequestOperationManager manager];
    NSDictionary *parameters = [userRegisterDataItem_ getDataItemEelements];
    NSString* urlString = [NSString stringWithFormat:@"%@%@", SERVER_IP, URL_PART_REGISTER_ONE_CLASS];
    [manager POST:urlString parameters:parameters success:^(AFHTTPRequestOperation *operation, id responseObject) {
        NSLog(@"JSON: %@", responseObject);
    } failure:^(AFHTTPRequestOperation *operation, NSError *error) {
        NSLog(@"Error: %@", error);
    }];
}

- (void)setEditable:(BOOL)bEditable
{
    classNameTextField_.enabled = bEditable;
    classTimeTextField_.enabled = bEditable;
    classTeachTextField_.enabled = bEditable;
    classDescriptionTextField_.enabled = bEditable;
    classMaxStudentNumberTextField_.enabled = bEditable;
}

- (void)updateUIFromData
{
    classNameTextField_.text = [notificationDataItem_ getClassName];
    classTimeTextField_.text = [notificationDataItem_ getClassTime];
    classTeachTextField_.text = [notificationDataItem_ getClassTeacher];
    classDescriptionTextField_.text = [notificationDataItem_ getClassDescription];
    classMaxStudentNumberTextField_.text = [notificationDataItem_ getClassMaxStudent];
}


@end
